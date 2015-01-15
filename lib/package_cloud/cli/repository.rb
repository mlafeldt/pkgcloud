module PackageCloud
  module CLI
    class Repository < Base
      option :private
      desc "create repo_name",
           "create repository called repo_name"
      def create(repo_name)
        config.read_or_create
        client = Client.new(config)

        print "Looking for repository at #{repo_name}... "
        repo = client.create_repo(repo_name, options[:private])
        print "success!\n".green
        puts "Your repository has been created at:"
        puts "    #{repo["url"]}"
      end

      desc "list",
           "list your repositories"
      def list
        repos = client.repositories
        if repos.length == 0
          puts "You have no repositories at the moment. Create one!"
        else
          puts "Your repositories:"
          puts ""
          repos.each_with_index do |repo, i|
            puts "  #{repo.fqname} (#{repo.private_human})"
            puts "  last push: #{repo.last_push_human} | packages: #{repo.package_count_human}"
            puts "" unless i == repos.length - 1
          end
        end
      end

      desc "install user/repo package_type",
           "install user/repo for package_type"
      def install(repo, package_type)
        if Process.uid != 0 && package_type != "gem"
          abort("You have to run install as root.".red)
        end

        print "Locating repository at #{repo}... "
        repo = client.repository(repo)
        print "success!\n"
        pid = fork { exec(repo.install_script(package_type)) }
        Process.waitpid(pid)
      end
    end
  end
end

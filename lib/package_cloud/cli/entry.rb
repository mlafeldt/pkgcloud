module PackageCloud
  module CLI
    class Entry < Base
      desc "repository SUBCMD ...ARGS", "manage repositories"
      subcommand "repository", Repository

      desc "distro SUBCMD ...ARGS", "manage repositories"
      subcommand "distro", Distro

      desc "master_token SUBCMD ...ARGS", "manage master tokens"
      subcommand "master_token", MasterToken

      desc "read_token SUBCMD ...ARGS", "manage read tokens"
      subcommand "read_token", ReadToken

      desc "yank user/repo[/distro/version] package_name",
           "yank package from user/repo [in dist/version]"
      def yank(repo_desc, package_name)
        ARGV.clear # otherwise gets explodes

        # strip os/dist
        repo_name = repo_desc.split("/")[0..1].join("/")
        dist = repo_desc.split("/")[2..3].join("/")

        if dist == "" && package_name =~ /\.gem$/
          dist = "gems"
        end

        print "Looking for repository at #{repo_name}... "
        repo = client.repository(repo_desc)
        print "success!\n"

        print "Attempting to yank package at #{repo_name}/#{dist}/#{package_name}..."
        packages = repo.yank(dist, package_name)
        puts "done!".green
      end

      desc "push user/repo[/distro/version] /path/to/packages",
           "push package(s) to repository (in distro/version if required)"
      option "skip-file-ext-validation", :type => :boolean
      option "yes", :type => :boolean
      def push(repo, package_file, *package_files)
        ARGV.clear # otherwise gets explodes
        package_files << package_file

        exts = package_files.map { |f| f.split(".").last }.uniq

        if package_files.length > 1 && exts.length > 1
          abort("You can't push multiple packages of different types at the same time.\nFor example, use *.deb to push all your debs at once.".red)
        end

        invalid_packages = package_files.select do |f|
          !["gem", "deb", "rpm", "dsc"].include?(f.split(".").last)
        end

        if !options.has_key?("skip-file-ext-validation") && invalid_packages.any?
          message = "I don't know how to push these packages:\n\n".red
          invalid_packages.each do |p|
            message << "  #{p}\n"
          end
          message << "\npackage_cloud only supports debs, gems, and rpms.".red
          abort(message)
        end

        if !options.has_key?("yes") && exts.first == "gem" && package_files.length > 1
          answer = get_valid("Are you sure you want to push #{package_files.length} packages? (y/n)") do |s|
            s == "y" || s == "n"
          end

          if answer != "y"
            abort("Aborting...".red)
          end
        end

        validator = Validator.new(client)
        dist_id   = validator.distribution_id(repo, package_files, exts.first)

        # strip os/dist
        repo = repo.split("/")[0..1].join("/")

        print "Looking for repository at #{repo}... "
        repo = client.repository(repo)
        print "success!\n"

        package_files.each do |f|
          files = nil
          ext = f.split(".").last

          if ext == "dsc"
            print "Checking source package #{f}... "
            files = parse_and_verify_dsc(repo, f, dist_id)
          end

          print "Pushing #{f}... "
          repo.create_package(f, dist_id, files, ext)
        end
      end

      desc "version",
           "print version information"
      def version
        puts "package_cloud CLI #{VERSION}\nSee https://packagecloud.io/docs#cli for more details."
      end

      private
        def parse_and_verify_dsc(repo, f, dist_id)
          files = repo.parse_dsc(f, dist_id)
          dirname = File.dirname(f)
          find_and_verify(dirname, files)
        end

        def find_and_verify(dir, files)
          file_paths = []
          files.each do |f|
            filepath = File.join(dir, f["filename"])
            if !File.exists?(filepath)
              print "Unable to find file name: #{f["filename"]} for source package: #{filepath}\n".red
              abort("Aborting...".red)
            end

            disk_size = File.stat(filepath).size
            if disk_size != f["size"]
              print "File #{f["filename"]} has size: #{disk_size}, expected: #{f["size"]}\n".red
              abort("Aborting...".red)
            end
            file_paths << filepath
          end
          file_paths
        end

        def get_valid(prompt)
          selection = ""
          times = 0
          until yield(selection)
            if times > 0
              puts "#{selection} is not a valid selection."
            end
            print "#{prompt}: "
            selection = ::Kernel.gets.chomp
            times += 1
          end

          selection
        end
    end
  end
end

module PackageCloud
  class Validator
    def initialize(client)
      @client = client
    end

    def distribution_id(repo, filenames, package_type)
      if distributions[package_type]
        _,_,dist_sel,ver_sel = repo.split("/")
        
        if dist_sel && ver_sel
          dist = distributions[package_type].detect { |d| d["index_name"] == dist_sel }

          if dist
            ver = dist["versions"].detect { |v| v["index_name"] == ver_sel || v["version_number"] == ver_sel }

            if ver
              ver["id"]
            else
              puts "#{ver_sel} isn't a valid version of #{dist["display_name"]}\n\n"
              select_dist(repo, filenames, package_type)
            end
          else
            puts "#{dist_sel} isn't a valid operating system.\n\n"
            select_dist(repo, filenames, package_type)
          end
        else
          puts "#{package_type} packages require an operating system and version to be selected.\n\n"
          select_dist(repo, filenames, package_type)
        end
      else
        nil
      end
    end

    private
      def distributions
        @distributions ||= @client.distributions
      end

      def select_from(list)
        selection = get_valid("0-#{list.length - 1}") do |s|
          s =~ /^\d+$/ && list[s.to_i]
        end

        list[selection.to_i]
      end

      def get_valid(prompt)
        selection = ""
        times = 0
        until yield(selection)
          if times > 0
            puts "#{selection} is not a valid selection."
          end
          print "\n #{prompt}: "
          selection = ::Kernel.gets.chomp
          times += 1
        end

        selection
      end

      def select_dist(repo, filenames, package_type)
        puts "If you don't see your OS or version here, send us an email at support@packagecloud.io:\n\n"
        all_distros = distributions[package_type]

        filtered_distros = all_distros.select {|dist| dist["index_name"] != "any"}
        
        filtered_distros.each_with_index do |dist, index|
          puts "\t#{index}. #{dist["display_name"]}"
        end

        distro = select_from(filtered_distros)

        puts "\nYou selected #{distro["display_name"]}. Select a version:\n\n"
        distro["versions"].each_with_index do |ver, index|
          puts "\t#{index}. #{ver["display_name"]} (#{ver["index_name"]})"
        end
        version = select_from(distro["versions"])

        repo = repo.split("/")[0..1].join("/")
        os_shortcut = "#{distro["index_name"]}/#{version["index_name"]}"
        shortcut = "#{repo}/#{os_shortcut}"
        if filenames.length > 1
          puts "\nPush #{filenames.length} packages to #{shortcut}? "
        else
          puts "\nPush #{filenames.first} to #{shortcut}? "
        end
        answer = get_valid("(y/n)") { |sel| sel == "y" || sel == "n" }

        if answer == "y"
          print "\nContinuing...".green 
          puts " Note that next time you can push directly to #{os_shortcut} by specifying #{shortcut} on the commandline."
          version["id"]
        else
          abort("Cancelled.")
        end
      end
  end
end

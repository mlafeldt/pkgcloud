module PackageCloud
  module CLI
    class Distro < Base
      desc "list package_type",
           "list available distros and versions for package_type"
      def list(package_type)
        distros = client.distributions[package_type]
        if distros
          puts "Listing distributions for #{package_type}:"
          distros.each do |distro|
            next if distro["index_name"] == "any"
            puts "\n    #{distro["display_name"]} (#{distro["index_name"]}):\n\n"
            distro["versions"].each do |ver|
              puts "        #{ver["display_name"]} (#{ver["index_name"]})"
            end
          end

          puts "\nIf you don't see your distribution or version here, email us at support@packagecloud.io."
        else
          puts "No distributions exist for #{package_type}.".red
          puts "That either means that we don't support #{package_type} or that it doesn't require a distribution/version."
          exit(1)
        end
      end
    end
  end
end

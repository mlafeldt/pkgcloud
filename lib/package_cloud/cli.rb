require "colorize"
require "thor"

module PackageCloud
  module CLI 
    autoload :Distro,      "package_cloud/cli/distro"
    autoload :Entry,       "package_cloud/cli/entry"
    autoload :MasterToken, "package_cloud/cli/master_token"
    autoload :ReadToken,   "package_cloud/cli/read_token"
    autoload :Repository,  "package_cloud/cli/repository"

    class Base < Thor
      class_option "config"
      class_option "url"

      private
        def config
          @config ||= begin
            ConfigFile.new(options[:config] || "~/.packagecloud",
                         options[:url] || "https://packagecloud.io").tap(&:read_or_create)
                      end
        end

        def client
          @client ||= Client.new(config)
        end
    end
  end
end

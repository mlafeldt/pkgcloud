require "highline/import"
require "json"
require "uri"
require "cgi"

module PackageCloud
  class ConfigFile
    attr_reader :token

    def initialize(filename = "~/.packagecloud", url = "https://packagecloud.io")
      abort_on_incorrect_locale
      @filename = File.expand_path(filename)
      @url = URI(url)
    end

    def read_or_create
      if ENV["PACKAGECLOUD_TOKEN"]
        @token = ENV["PACKAGECLOUD_TOKEN"]
        @url   = URI(ENV["PACKAGECLOUD_URL"]) if ENV["PACKAGECLOUD_URL"]
        assert_reasonable_gem_version
      elsif File.exist?(@filename)
        attrs = JSON.parse(File.read(@filename))
        @token = attrs["token"] if attrs.has_key?("token")
        @url   = URI(attrs["url"]) if attrs.has_key?("url")
        assert_reasonable_gem_version
      else
        puts "No config file exists at #{@filename}. Login to create one."

        @token = login_from_console
        write
      end
    end

    def url
      @url ||= URI("https://packagecloud.io")
    end

    def base_url(username = token, password = "")
      u = url.dup
      u.user = CGI.escape(username)
      u.password = CGI.escape(password)
      u.to_s
    end

    private

      def abort_on_incorrect_locale
        default_encoding = Encoding.default_external
        if default_encoding != Encoding::UTF_8
          message = "It appears your locale is set to #{default_encoding}, " +
              "please use UTF-8 instead. Run the 'locale' command to get" +
              "more info or e-mail support@packagecloud.io for help"
          abort(message.red)
        end
      end

      def login_from_console
        e     = ask("Email:")
        p     = ask("Password:") { |q| q.echo = false }

        begin
          PackageCloud::Auth.get_token(base_url(e, p))
        rescue RestClient::Unauthorized => e
          puts "Sorry, but we couldn't find you. Give it another try."
          login_from_console
        end
      end

      def write
        print "Got your token. Writing a config file to #{@filename}... "
        attrs = {url => url.to_s, :token => @token}
        File.open(@filename, "w", 0600) { |f| f << JSON.dump(attrs); f << "\r\n" }
        puts "success!"
      end

      def assert_reasonable_gem_version
        gem_version = Client.new(self).gem_version
        if gem_version["minor"] != MINOR_VERSION
          abort("This gem is out of date. Please update it!".red)
        elsif gem_version["patch"] > PATCH_VERSION
          $stderr << "[WARNING] There's a newer version of the package_cloud gem. Install it when you get a chance!\n"
        end
      end
  end
end

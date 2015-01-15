require "json"
require "rest_client"

module PackageCloud
  class Client
    def initialize(config)
      @config = config
    end

    def repositories
      base_url = @config.base_url
      begin
        attrs = JSON.parse(RestClient.get("#{base_url}/api/v1/repos.json"))
        attrs.map { |a| Repository.new(a, @config) }
      rescue RestClient::ResourceNotFound => e
        print "failed!\n".red
        exit(127)
      end
    end

    def repository(name)
      base_url = @config.base_url
      user, repo = name.split("/")
      begin
        attrs = JSON.parse(RestClient.get("#{base_url}/api/v1/repos/#{user}/#{repo}.json"))
        if attrs["error"] == "not_found"
          print "failed... Repository #{user}/#{repo} not found!\n".red
          exit(127)
        end

        Repository.new(attrs, @config)
      rescue RestClient::ResourceNotFound => e
        print "failed!\n".red
        exit(127)
      end
    end

    def create_repo(name, priv)
      base_url = @config.base_url
      begin
        JSON.parse(RestClient.post("#{base_url}/api/v1/repos.json", :repository => {:name => name, :private => priv == "private" ? "1" : "0"}))
      rescue RestClient::UnprocessableEntity => e
        print "error!\n".red
        json = JSON.parse(e.response)
        json.each do |k,v|
          puts "\n\t#{k}: #{v.join(", ")}\n"
        end
        puts ""
        exit(1)
      end
    end

    def distributions
      base_url = @config.base_url
      begin
        JSON.parse(RestClient.get("#{base_url}/api/v1/distributions.json"))
      rescue RestClient::ResourceNotFound => e
        print "failed!\n".red
        exit(127)
      end
    end

    def gem_version
      base_url = @config.base_url
      begin
        JSON.parse(RestClient.get("#{base_url}/api/v1/gem_version"))
      rescue RestClient::ResourceNotFound => e
        print "failed!\n".red
        exit(127)
      end
    end
  end
end

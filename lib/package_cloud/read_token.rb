module PackageCloud
  class ReadToken < Object
    def initialize(config, attrs)
      @config = config
      @attrs = attrs
    end

    def destroy(master_token_path, read_token_id)
      base_url = @config.base_url
      url = "#{base_url}#{master_token_path}/read_tokens/#{read_token_id}"
      RestClient.delete(url)
    end
  end
end

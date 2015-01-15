module PackageCloud
  class MasterToken < Object
    def initialize(attrs, config)
      @attrs = attrs
      @config = config
    end

    def read_tokens
      @attrs["read_tokens"].map do |read_token|
        ReadToken.new(@config, read_token)
      end
    end

    def destroy
      url = @config.base_url + paths["self"]
      RestClient.delete(url)
    end
  end
end

require "rest_client"
require "json"

module PackageCloud
  module Auth
    class << self
      def get_token(url)
        JSON.parse(RestClient.get("#{url}/api/v1/token.json"))["token"]
      end
    end
  end
end

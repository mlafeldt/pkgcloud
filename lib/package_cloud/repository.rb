module PackageCloud
  class Repository < Object
    def initialize(attrs, config)
      @attrs = attrs
      @config = config
    end

    def parse_dsc(dsc_path, dist_id)
      file_data = File.new(dsc_path, 'rb')
      base_url = @config.base_url
      url = base_url + paths["package_contents"]
      begin
        resp = RestClient::Request.execute(:method => 'post',
                                           :url => url,
                                           :timeout => -1,
                                           :payload => { :package => {:package_file      => file_data,
                                                                      :distro_version_id => dist_id}})
        resp = JSON.parse(resp)
        print "success!\n"
        resp["files"]
      rescue RestClient::UnprocessableEntity => e
        print "error:\n".red
        json = JSON.parse(e.response)
        json.each do |k,v|
          puts "\n\t#{k}: #{v.join(", ")}\n"
        end
        puts ""
        exit(1)
      end
    end

    def create_package(file_path, dist_id, files=nil, filetype=nil)
      file_data = File.new(file_path, 'rb')
      base_url = @config.base_url
      url = base_url + paths["create_package"]
      params = { :package_file => file_data,
                 :distro_version_id => dist_id }

      if filetype == "dsc"
        file_ios = files.inject([]) do |memo, f|
          memo << File.new(f, 'rb')
        end
        params.merge!({:source_files => file_ios})
      end

      begin
        RestClient::Request.execute(:method => 'post',
                                    :url => url,
                                    :timeout => -1,
                                    :payload => { :package =>  params })
        print "success!\n".green
      rescue RestClient::UnprocessableEntity => e
        print "error:\n".red
        json = JSON.parse(e.response)
        json.each do |k,v|
          puts "\n\t#{k}: #{v.join(", ")}\n"
        end
        puts ""
        exit(1)
      end
    end

    def install_script(type)
      url = urls["install_script"].gsub(/:package_type/, type)
      RestClient.get(url)
    end

    def master_tokens
      url = @config.base_url + paths["master_tokens"]
      attrs = JSON.parse(RestClient.get(url))
      attrs.map { |a| MasterToken.new(a, @config) }
    end

    def create_master_token(name)
      url = @config.base_url + paths["create_master_token"]
      begin
        RestClient.post(url, :master_token => {:name => name})
      rescue RestClient::UnprocessableEntity => e
        print "error:\n".red
        json = JSON.parse(e.response)
        json.each do |k,v|
          puts "\n\t#{k}: #{v.join(", ")}\n"
        end
        puts ""
        exit(1)
      end
    end

    def yank(dist, package_name)
      url = @config.base_url + paths["self"] + "/" + [dist, package_name].compact.join("/")
      begin
        RestClient.delete(url)
      rescue RestClient::ResourceNotFound => e
        print "error:\n".red
        json = JSON.parse(e.response)
        json.each do |k,v|
          puts "\n\t#{k}: #{v.join(", ")}\n"
        end
        puts ""
        exit(1)
      end
    end

    def private_human
      send(:private) ? "private".red : "public".green
    end
  end
end

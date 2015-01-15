module PackageCloud
  module CLI
    class ReadToken < Base
      desc "destroy user/repository mastertoken/readtoken",
           "revokes read token associated to a master token"
      def destroy(repo_name, master_and_read_token)
        print "Looking for repository at #{repo_name}... "
        repo = client.repository(repo_name)

        given_master_token, given_read_token = master_and_read_token.split("/")

        if given_master_token.nil? || given_read_token.nil?
          print "invalid master token and/or read token!\n".red
          exit(127)
        end

        master_token = repo.master_tokens.detect { |t| t.name == given_master_token }

        if master_token.nil?
          print "couldn't find master token named #{given_master_token}\n".red
          exit(127)
        end

        read_token = master_token.read_tokens.detect { |t| t.name == given_read_token }

        if read_token.nil?
          print "couldn't find read token named #{given_read_token} for #{given_master_token}\n".red
          exit(127)
        end

        master_token_path = master_token.paths["self"]

        read_token.destroy(master_token_path, read_token.id)

        print "success!\n"
      end
    end
  end
end
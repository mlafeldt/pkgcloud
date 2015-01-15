module PackageCloud
  class Object
    def respond_to?(method, include_priv=false)
      @attrs.has_key?(method.to_s) || super
    end

    def method_missing(method, *args, &block)
      if @attrs.has_key?(method.to_s)
        @attrs[method.to_s]
      else
        super
      end
    end
  end
end

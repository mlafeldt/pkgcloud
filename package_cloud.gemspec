# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'package_cloud/version'

Gem::Specification.new do |spec|
  spec.name          = "package_cloud"
  spec.version       = PackageCloud::VERSION
  spec.authors       = ["Joe Damato"]
  spec.email         = ["support@packagecloud.io"]
  spec.description   = %q{https://packagecloud.io}
  spec.summary       = %q{https://packagecloud.io}
  spec.homepage      = "https://packagecloud.io"
  spec.license       = "MIT"

  spec.files         = `git ls-files`.split($/)
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.require_paths = ["lib"]

  spec.add_runtime_dependency "thor", "0.18.1"
  spec.add_runtime_dependency "highline", "1.6.20"
  spec.add_runtime_dependency "rest-client", "1.6.7"
  spec.add_runtime_dependency "json_pure", "1.8.1"
  spec.add_runtime_dependency "colorize", "0.6.0"

  spec.add_development_dependency "bundler", "~> 1.3"
  spec.add_development_dependency "rake"
end

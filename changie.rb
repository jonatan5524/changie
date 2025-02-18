# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Changie < Formula
  desc "Automated changelog tool for preparing releases with lots of customization options."
  homepage "https://changie.dev"
  version "1.7.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/miniscruff/changie/releases/download/v1.7.0/changie_1.7.0_darwin_amd64.tar.gz"
      sha256 "f62f5a50c3dac5bc703a3c3013411be6785e592b9d5a4b2ce2c64232af4aeeb2"

      def install
        bin.install "changie"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/miniscruff/changie/releases/download/v1.7.0/changie_1.7.0_darwin_arm64.tar.gz"
      sha256 "7fdcad929a8a4d8de983061f6dce9689553cce2702d35e01dcbe96791a06ec34"

      def install
        bin.install "changie"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/miniscruff/changie/releases/download/v1.7.0/changie_1.7.0_linux_amd64.tar.gz"
      sha256 "92519b7439bbfa544070c36514ee23fc80f0717b19d38463fbc1a479aba5651d"

      def install
        bin.install "changie"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/miniscruff/changie/releases/download/v1.7.0/changie_1.7.0_linux_arm64.tar.gz"
      sha256 "97ba5ef818ccb7fc4717cf4cd90e3cf5fc11cc0d4e041d00cb566d367f0b1ce5"

      def install
        bin.install "changie"
      end
    end
  end
end

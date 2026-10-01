# Homebrew formula for Carrel. Filled in by packaging/fill.sh; see packaging/README.md.
class Carrel < Formula
  desc "Practise data structures and algorithms locally in Java or C++"
  homepage "https://github.com/saurabhraghuvanshii/carrel"
  version "{{version}}"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/saurabhraghuvanshii/carrel/releases/download/v{{version}}/carrel_darwin_arm64.tar.gz"
      sha256 "{{sha_darwin_arm64}}"
    end
    on_intel do
      url "https://github.com/saurabhraghuvanshii/carrel/releases/download/v{{version}}/carrel_darwin_amd64.tar.gz"
      sha256 "{{sha_darwin_amd64}}"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/saurabhraghuvanshii/carrel/releases/download/v{{version}}/carrel_linux_arm64.tar.gz"
      sha256 "{{sha_linux_arm64}}"
    end
    on_intel do
      url "https://github.com/saurabhraghuvanshii/carrel/releases/download/v{{version}}/carrel_linux_amd64.tar.gz"
      sha256 "{{sha_linux_amd64}}"
    end
  end

  def install
    bin.install "carrel"
  end

  def caveats
    <<~EOS
      Carrel runs your code with the compilers already on your computer.
      Run `carrel doctor` to see whether Java and g++ are installed.
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/carrel version")
  end
end

class Usage < Formula
  desc "Serve AI coding tool token usage over HTTP"
  homepage "https://github.com/notzorexlol/usage"
  url "https://github.com/notzorexlol/usage/archive/refs/tags/v0.1.0.tar.gz"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"usage", "./cmd/usage"
  end

  service do
    run [opt_bin/"usage", "-port", "7679"]
    keep_alive true
  end

  test do
    system bin/"usage", "-h"
  end
end

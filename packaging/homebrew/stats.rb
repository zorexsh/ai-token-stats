class Stats < Formula
  desc "Serve AI coding tool token stats over HTTP"
  homepage "https://github.com/notzorexlol/stats"
  url "https://github.com/notzorexlol/stats/archive/refs/tags/v0.1.0.tar.gz"
  license "MIT"
  head "https://github.com/notzorexlol/stats.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/stats"
  end

  service do
    run [opt_bin/"stats", "-port", "7679"]
    keep_alive true
    working_dir HOMEBREW_PREFIX
  end

  test do
    assert_match "port to listen on", shell_output("#{bin}/stats -h 2>&1")
  end
end
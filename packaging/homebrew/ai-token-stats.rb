class AiTokenStats < Formula
  desc "Local AI token stats for a portfolio page"
  homepage "https://github.com/zorexsh/ai-token-stats"
  url "https://github.com/zorexsh/ai-token-stats/archive/refs/tags/v0.1.0.tar.gz"
  license "MIT"
  head "https://github.com/zorexsh/ai-token-stats.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/ai-token-stats"
  end

  service do
    run [opt_bin/"ai-token-stats", "-port", "7679"]
    keep_alive true
    working_dir HOMEBREW_PREFIX
  end

  test do
    assert_match "port to listen on", shell_output("#{bin}/ai-token-stats -h 2>&1")
  end
end

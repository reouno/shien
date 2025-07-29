class Shien < Formula
  desc "Background daemon application to support knowledge workers"
  homepage "https://github.com/yourusername/shien"
  url "https://github.com/yourusername/shien/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"
  head "https://github.com/yourusername/shien.git", branch: "main"

  depends_on "go" => :build

  def install
    system "make", "build-all"
    bin.install "shien"
    bin.install "shienctl"
  end

  service do
    run [opt_bin/"shien"]
    keep_alive true
    log_path var/"log/shien.log"
    error_log_path var/"log/shien.err.log"
  end

  test do
    # Test shienctl version
    assert_match "shienctl", shell_output("#{bin}/shienctl --help 2>&1")
    
    # Test daemon startup (in background)
    pid = fork { exec bin/"shien" }
    sleep 2
    
    # Test ping command
    output = shell_output("#{bin}/shienctl ping 2>&1")
    assert_match "pong", output
    
    # Clean up
    Process.kill("TERM", pid)
    Process.wait(pid)
  end
end
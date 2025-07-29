class Shien < Formula
  desc "Background daemon application to support knowledge workers"
  homepage "https://github.com/yourusername/shien"
  url "https://github.com/yourusername/shien/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "実際のSHA256ハッシュ値をここに入れる"
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
    assert_match "shienctl", shell_output("#{bin}/shienctl --help 2>&1")
  end
end
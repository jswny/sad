class Sad < Formula
  desc "Simple app deployment based on SSH and Docker."
  homepage "https://github.com/jswny/sad"
  url "https://github.com/jswny/sad/archive/v3.0.1.tar.gz"
  sha256 "92820e37c7ae43b8d469359bd7c40a049a8ab6697df433c140bd82f59c450500"
  license ""

  depends_on "openssh"
  depends_on "go" => :build

  def install
    system "go", "build", "-o=sad", "cmd/sad/main.go"
    bin.install "sad"
  end

  test do
    assert_match "Usage of", shell_output("#{bin}/sad -help", 2).strip
  end
end

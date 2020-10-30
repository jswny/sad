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
    # `test do` will create, run in and delete a temporary directory.
    #
    # This test will fail and we won't accept that! For Homebrew/homebrew-core
    # this will need to be a test that verifies the functionality of the
    # software. Run the test with `brew test sad`. Options passed
    # to `brew install` such as `--HEAD` also need to be provided to `brew test`.
    #
    # The installed folder is not in the path, so use the entire path to any
    # executables being tested: `system "#{bin}/program", "do", "something"`.
    # assert_equal "system "#{bin}/sad", "-help", "something"
    assert_match "Usage of", shell_output("#{bin}/sad -help", 2).strip
  end
end

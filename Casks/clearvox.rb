cask "clearvox" do
  version "1.0.0"
  sha256 :no_check  # Update with actual SHA256 after first release

  url "https://github.com/errakhaoui/noise-canceling/releases/download/v#{version}/ClearVox-Installer.dmg"
  name "ClearVox"
  desc "Real-time noise cancellation for macOS with virtual microphone support"
  homepage "https://github.com/errakhaoui/noise-canceling"

  livecheck do
    url :url
    strategy :github_latest
  end

  # Extract PKG from DMG and install
  installer manual: "ClearVox-Installer.pkg"

  # Alternative: Direct app installation (if switching back to app-in-dmg)
  # app "ClearVox.app"

  zap trash: [
    "~/Library/Application Support/ClearVox",
    "~/Library/Preferences/com.errakhaoui.clearvox.plist",
    "~/Library/Caches/com.errakhaoui.clearvox",
  ]

  caveats <<~EOS
    ClearVox requires microphone access to function.

    On first launch:
    1. Grant microphone permission when prompted
    2. (Optional) Install BlackHole for virtual microphone support
    3. Configure audio devices and click Start

    For virtual microphone in video calls:
      brew install blackhole-2ch

    Documentation: https://github.com/errakhaoui/noise-canceling
  EOS
end

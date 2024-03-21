# 💽 smv _(split my video)_

A simple to use tool, that splits youtube videos based off silence using
[khdai's youtube library](https://github.com/kkdai/youtube) and [ffmpeg](https://ffmpeg.org/).

#### 🤔 Who is this for?

If you like music discographies from youtube that aren't readily available on mainstream music streaming services
and like keeping music locally, then this is for you!

Manually splitting up songs from videos using something like audacity is a pain and a chore,
and so I made a tool to do that for me! And for you too, if you stumbled on this repo and think it'd be useful to you that is.

## ⭐️ Features

- Split audio from youtube videos based on silence.

  _That's all its useful for, for now at least._

## 💥 Installation

### Using go get

Ensure you have Go 1.22.0 or later.

```go
go get github.com/vague2k/smv
```

### Homebrew (macOS)

```
brew install smv
```

### From Source

```
git clone https://github.com/vague2k/smv.git
cd smv
go run main.go
```

You can find

## 📋 Contributing

Issues and PR's are always welcome and highly encouraged! as this is my first plugin. I would love to learn more.

## License

[MIT](https://choosealicense.com/licenses/mit/)

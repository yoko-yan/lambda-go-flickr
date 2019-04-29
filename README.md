# lambda-go-flickr

AWS Lambda のプログラム。言語はGo
位置情報（緯度経度）を指定して、Flickrの画像URLを返す関数

動作には、FlickrのAPIキーが必要です。
APIキーは、ココで作成できます。
https://www.flickr.com/services/apps/create/

使用しているFlickrのAPIのリファレンスはココ
https://www.flickr.com/services/api/flickr.photos.search.html

Flickrの画像のURLの組み立てについては、ココ
https://www.flickr.com/services/api/misc.urls.html

FlickrのAPIを試す場合は、ココ
https://www.flickr.com/services/api/explore/flickr.photos.search

AWSにアップロードする際は、以下のコマンドでビルド後、ZIPファイルを作って、アップロードしてください。
$ GOOS=linux GOARCH=amd64 go build -o hello
$ zip hello.zip hello


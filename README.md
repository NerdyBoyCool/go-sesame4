# go-sesame4
Github Actionを使用して、[Sesame4](https://jp.candyhouse.co/products/sesame4)の施錠を自動化するGolang Script

## Motivation
オートロック機能もあるが、「鍵を持ったかどうか？」という新しい悩みが出るので使いたくない。
やりたいことは、深夜の鍵の閉め忘れを防ぐために、一定時間になったら施錠したいというだけ。アプリから施錠はできるが、そもそも「閉めたかどうか？」を気にしたくないというのがモチベーション。

## Sesame Web API
https://doc.candyhouse.co/ja/SesameAPI

## API KEYS
招待用のQRを発行し、そのQRコードを下記のURLからデコードして、必要な情報得る必要がある。

- SESAME_API_KEY: Public Key
- SESAME_SECRET_KEY: Secret Key
- SESAME_UUID: UUID

https://sesame-qr-reader.vercel.app/

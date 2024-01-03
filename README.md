<h1>COCOLOTalk / 対人恐怖症の方ためのSNS</h1>
    <p>
        フロントエンドのリポジトリは<a href="https://github.com/KeitaShimura/SadAppClient" target="_blank">こちら</a>
    </p>
    <img width="1800" alt="スクリーンショット 2023-12-08 2 56 59" src="https://github.com/KeitaShimura/SadAppAPI/assets/124238548/b159463e-247d-4f25-8bc5-84562412837d">
    <img width="1800" alt="スクリーンショット 2023-12-08 2 56 44" src="https://github.com/KeitaShimura/SadAppAPI/assets/124238548/d96ebfca-0aef-4fb7-ad7b-1da56136f167">

</p>
<h2>サービス概要</h2>
<p>COCOLOTalkは「対人恐怖症、社交不安障害（SAD）の方のお悩みを解決したい！」という想いから作られた、無料のSNSです。</p>
<p>自分と同じお悩みを持つ方に悩みを打ち明けることができます。</p>

<h2>メイン機能の使い方</h2>
![画面収録 2024-01-04 1 49 05 (1)](https://github.com/KeitaShimura/SadAppAPI/assets/124238548/1785d91e-8ece-4fb9-9909-3acce63bb6f9)


<h2>使用技術一覧</h2>

- バックエンド: Go / Fiber

  - コード解析: golangci-lint
  - フォーマッター: gofmt
  - テストパッケージ: testing

- フロントエンド: JavaScript / React

  - コード解析: ESLint
  - フォーマッター: Prettier
  - テストフレームワーク: React Testing Library / Jest
  - 主要パッケージ: Axios / Font Awesome / React Bootstrap / React Toastify

- CI / CD: GitHub Actions
- 環境構築: Docker / Docker Compose
- インフラ: Render / Nginx / Vercel

<h2>機能一覧</h2>
    <ul>
        <li>ログイン</li>
        <li>新規登録</li>
        <li>ログアウト</li>
        <li>投稿、投稿一覧、削除</li>
        <li>イベント、イベント一覧、削除</li>
        <li>ユーザー一覧、フォロー、フォロワー一覧</li>
        <li>プロフィール、プロフィール編集</li>
    </ul>

<h2>画面</h2>
    <ul>
        <li>トースト表示</li>
        <li>モーダル画面(各画面の詳細は下記の画面遷移図参照)</li>
        <li>404エラーのカスタム画面</li>
        <li>レスポンシブデザイン</li>
    </ul>

<h2>ER図</h2>
<img width="888" alt="スクリーンショット 2023-12-11 2 53 06" src="https://github.com/KeitaShimura/SadAppAPI/assets/124238548/5240faf4-c484-4969-af10-01ad3fe44d48" /></br>
<a href="https://dbdiagram.io/d/64600a51dca9fb07c40853cc" target="_blank">ER図</a>

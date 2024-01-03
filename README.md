<h1>COCOLOTalk / 対人恐怖症の方ためのSNS</h1>
    <p>
        フロントエンドのリポジトリは<a href="https://github.com/KeitaShimura/SadAppClient" target="_blank">こちら</a>
    </p>
    <img width="1800" alt="スクリーンショット 2023-12-08 2 56 59" src="https://github.com/KeitaShimura/SadAppAPI/assets/124238548/b159463e-247d-4f25-8bc5-84562412837d">
    <img width="1800" alt="スクリーンショット 2023-12-08 2 56 44" src="https://github.com/KeitaShimura/SadAppAPI/assets/124238548/d96ebfca-0aef-4fb7-ad7b-1da56136f167">

<h2>技術スタック</h2>
<p>
        <a href="https://golang.org/" target="_blank">
            <img src="https://img.shields.io/badge/-Go-00ADD8.svg?logo=go&style=flat-square&logoColor=white" alt="Go Badge">
        </a>
        <a href="https://gofiber.io/" target="_blank">
            <img src="https://img.shields.io/badge/-Go_Fiber-88C0D0.svg?logo=gofiber&style=flat-square&logoColor=white" alt="Go-Fiber Badge">
        </a>
        <a href="https://reactjs.org/" target="_blank">
            <img src="https://img.shields.io/badge/-React-61DAFB.svg?logo=react&style=flat-square&logoColor=white" alt="React Badge">
        </a>
        <a href="https://developer.mozilla.org/en-US/docs/Web/JavaScript">
            <img src="https://img.shields.io/badge/-JavaScript-F7DF1E.svg?logo=javascript&style=flat-square&logoColor=black" alt="JavaScript Badge">
        </a>
        <a href="https://www.docker.com/" target="_blank">
            <img src="https://img.shields.io/badge/-Docker-2496ED.svg?logo=docker&style=flat-square&logoColor=white" alt="Docker Badge">
        </a>
        <a href="https://www.atlassian.com/continuous-delivery/principles/continuous-integration-vs-delivery-vs-deployment" target="_blank">
            <img src="https://img.shields.io/badge/-CI%2FCD-2088FF.svg?style=flat-square" alt="CI/CD Badge">
        </a>
</p>

<h2>サービス概要</h2>
<p>COCOLOTalkは「対人恐怖症、社交不安障害（SAD）の方のお悩みを解決したい！」という想いから作られた、無料のSNSです。</p>
<p>自分と同じお悩みを持つ方に悩みを打ち明けることができます。</p>

<h2>メイン機能の使い方</h2>
    <table border="1">
        <tr>
            <th>アカウント登録</th>
            <th>ヘッダ2</th>
            <th>ヘッダ3</th>
        </tr>
        <tr>
            <td>　![画面収録 2024-01-04 1 49 05](https://github.com/KeitaShimura/SadAppAPI/assets/124238548/d20a20ac-a093-4b46-8b50-6595beeb2c76)
　</td>
            <td>セル1-2</td>
            <td>セル1-3</td>
        </tr>
    </table>
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

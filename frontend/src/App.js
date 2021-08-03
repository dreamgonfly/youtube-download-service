import React, { Fragment, useState } from 'react';
import api from './api';

function App() {
  const [enteredVideoId, setEnteredVideoId] = useState('');
  const [data, setData] = useState({ Thumbnail: "", Formats: [] })

  const videoIdChangeHandler = (event) => {
    setEnteredVideoId(event.target.value);
  };

  const previewClickHandler = (event) => {
    api.preview(enteredVideoId).then((res) => {
      console.log(res)
      setData(res.data)

    })
  }

  // const data = {
  //   Thumbnail: "https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/%5B%EA%B8%B0%EC%83%9D%EC%B6%A9%5D%2030%EC%B4%88%20%EC%98%88%EA%B3%A0.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256&X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210803%2Fauto%2Fstorage%2Fgoog4_request&X-Goog-Date=20210803T004718Z&X-Goog-Expires=899&X-Goog-Signature=9559009897db3c9d382ed73cbbfc3d349f00b798b4ab6e0980d4acabaca09f87757721933b97368d69c1f6076aa1a01f8ea7ff9ce3d3aa0df8d700bccace51f2487f168d238ab1b48b5986172bc140c3020251b22b6ea5a3fac6f588e6f5cbf192fa4a5178d9879dfe067622735350bbd14a376f5847dcdfe13ae54799222e78295372843e335f9d55172c7b98cfe71dd3cc79d7bbf53a7c908be9b4bb7b589b93337cfa0e01c08aca5dd33f025833168f1fc0ee6e48c411ed68b968d1ff2775a386a3a6b8b99962e13721785c133b2e300b1e957160d5a1c01d1f54f6719e298bd316af0fcab50c110c37ef26a878cb2ac5723e7bfcf186462c0f428e878b0b&X-Goog-SignedHeaders=host",
  //   Formats: [
  //     {
  //       Filesize: 1348634,
  //       FormatId: "18",
  //       FormatNote: "360p",
  //       Ext: "mp4"
  //     },
  //     {
  //       Filesize: 2059470,
  //       FormatId: "22",
  //       FormatNote: "720p",
  //       Ext: "mp4"
  //     }
  //   ]
  // }

  let thumbnailElement = <div></div>
  if (data['Thumbnail'] !== "") {
    thumbnailElement = <Fragment><h3>Thumbnail</h3><img src={data.Thumbnail} alt="thumbnail" /></Fragment>
  }

  return (
    <div>
      <h1>Youtube Download Service</h1>
      <input type="text" onChange={videoIdChangeHandler} />
      <button onClick={previewClickHandler}>Preview</button>
      <br />
      {data.Formats.map((item) => {
        return <Fragment>
          <br />
          <a href={api.composeDownloadLink(enteredVideoId, item.FormatId)}><button>Download {item.FormatNote}</button></a>
          <br />
        </Fragment>
      })}
      {thumbnailElement}
    </div>
  );
}

export default App;

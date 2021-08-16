import React, { useState } from 'react'
import styles from './OutputAction.module.css'
import api from '../../api'
import OpenMessage from './OpenMessage'
import OptionsBlock from './OptionsBlock'

const OutputAction = (props) => {
    // Existance of props.data is guaranteed
    const [formatId, setFormatId] = useState(props.data.Formats[0].FormatId)
    const [optionsOpen, setOptionsOpen] = useState(false)
    const [loading, setLoading] = useState(false)

    const format = props.data.Formats.find(element => element.FormatId === formatId)
    const fileName = props.data.Name + "." + format.Ext
    const fileSizeInMB = Math.round(format.Filesize / 1000 / 1000 * 10) / 10

    const downloadLink = api.composeDownloadLink(props.videoId, format.FormatId, fileName)

    const choiceClickHandler = () => {
        setOptionsOpen((prev) => { return !prev })
    }

    const playClickHandler = () => {
        setLoading(true)
        api.play(props.videoId, format.FormatId, fileName).then((res) => {
            setLoading(false)
            // Secure solution. https://stackoverflow.com/questions/45046030/maintaining-href-open-in-new-tab-with-an-onclick-handler-in-react
            const newWindow = window.open(res.data.URL, '_blank', 'noopener,noreferrer')
            if (newWindow) newWindow.opener = null
        })
    }


    return (
        <div className={styles["output-download"]}>
            <div className={styles["output-format-choice"]}>
                <div className={styles["format-desc"]}>{format.Ext} {format.FormatNote} {fileSizeInMB} MB</div>
                <div className={styles["format-choice-arrow"]} onClick={choiceClickHandler}>
                    <i className={`fas fa-caret-down ${styles["arrow-down"]}`}></i>
                </div>
            </div>
            {optionsOpen ? <OptionsBlock data={props.data} formatId={formatId} onSetFormatId={setFormatId} onSetOptionsOpen={setOptionsOpen} /> : <div></div>}
            <div>
                <div className={styles["output-action"]}>
                    <a className={styles['output-link']} href={downloadLink}>
                        <button className={styles["output-download-button"]}>Download</button>
                    </a>
                </div>
                <div className={styles["output-action"]}>
                    <button className={styles["output-view"]} onClick={playClickHandler}>
                        <OpenMessage isLoading={loading} />
                    </button>
                </div>
            </div>
        </div>
    )
}

export default OutputAction
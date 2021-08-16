import React, { useState } from 'react'
import styles from './OutputAction.module.css'
import api from '../../api'
import OpenMessage from './OpenMessage'
import OptionsBlock from './OptionsBlock'

const OutputAction = (props) => {
    let [selectedIndex, setSelectedIndex] = useState(undefined)
    const [optionOpen, setOptionOpen] = useState(false)
    const [loading, setLoading] = useState(false)

    if (props.data.Formats.length > 0 && selectedIndex === undefined) {
        selectedIndex = props.data.Formats[0].FormatId
    }

    const viewClickHandlerConstructur = (id, format, filename) => {
        const viewClickHandler = () => {
            setLoading(true)
            api.save(id, format, filename).then((res) => {
                setLoading(false)
                // Secure solution. https://stackoverflow.com/questions/45046030/maintaining-href-open-in-new-tab-with-an-onclick-handler-in-react
                const newWindow = window.open(res.data.URL, '_blank', 'noopener,noreferrer')
                if (newWindow) newWindow.opener = null
            })
        }
        return viewClickHandler
    }

    const choiceClickHandler = () => {
        setOptionOpen((prev) => { return !prev })
    }

    return (
        <div className={styles["output-download"]}>
            <div className={styles["output-format-choice"]}>
                <div className={styles["format-desc"]}>{props.data.Formats.find(element => element.FormatId === selectedIndex).Ext} {props.data.Formats.find(element => element.FormatId === selectedIndex).FormatNote} {Math.round(props.data.Formats.find(element => element.FormatId === selectedIndex).Filesize / 1000 / 1000 * 10) / 10} MB</div>
                <div className={styles["format-choice-arrow"]} onClick={choiceClickHandler}><i className={`fas fa-caret-down ${styles["arrow-down"]}`}></i></div>
            </div>
            {optionOpen ? <OptionsBlock data={props.data} selectedIndex={selectedIndex} onSetSelectedIndex={setSelectedIndex} setOptionOpen={setOptionOpen} /> : <div></div>}
            <div>
                <div className={styles["output-action"]}>
                    <a className={styles['output-link']} href={api.composeDownloadLink(props.videoId, props.data.Formats.find(element => element.FormatId === selectedIndex).FormatId, props.data.Name + "." + props.data.Formats.find(element => element.FormatId === selectedIndex).Ext)}><button className={styles["output-download-button"]}>Download</button></a>
                </div>
                <div className={styles["output-action"]}>
                    <button className={styles["output-view"]} onClick={viewClickHandlerConstructur(props.videoId, props.data.Formats.find(element => element.FormatId === selectedIndex).FormatId, props.data.Name + "." + props.data.Formats.find(element => element.FormatId === selectedIndex).Ext)}>
                        <OpenMessage isLoading={loading} />
                    </button>
                </div>
            </div>
        </div>
    )
}

export default OutputAction
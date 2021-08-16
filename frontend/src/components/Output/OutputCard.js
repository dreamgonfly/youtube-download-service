import React from 'react'
import OutputAction from './OutputAction'
import OutputContent from './OutputContent'
import styles from './OutputCard.module.css'

const OutputCard = (props) => {

    if (props.data.Formats.length === 0) {
        return <div></div>
    }

    return (
        <div className={styles["output-card"]}>
            <OutputContent data={props.data} />
            <OutputAction videoId={props.videoId} data={props.data} />
        </div>
    )
}

export default OutputCard
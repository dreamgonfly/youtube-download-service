import React from 'react'
import styles from './OutputContent.module.css'

const OutputContent = (props) => {

    const minutes = Math.floor(props.data.DurationSecond / 60)
    const seconds = props.data.DurationSecond - minutes * 60

    return (
        <div className={styles["output-content"]}>
            <div className={styles["output-thumbnail"]}>
                <img className={styles["img"]} src={props.data.Thumbnail} alt="thumbnail" />
            </div>
            <div className={styles["output-title"]}>{props.data.Title}</div>
            <div className={styles["output-length"]}>{minutes}:{seconds}</div>
        </div>
    )
}

export default OutputContent


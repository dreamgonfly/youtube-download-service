import React from 'react'
import OutputAction from './OutputAction'
import OutputContent from './OutputContent'
import styles from './OutputCard.module.css'
import Loader from "react-loader-spinner"

const OutputCard = (props) => {

    if (props.waiting) {
        return (
            <div className={styles["output-spin"]}>
                <Loader
                    type="TailSpin"
                    color="#00BFFF"
                    height={100}
                    width={100}
                />
            </div>
        )
    }

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
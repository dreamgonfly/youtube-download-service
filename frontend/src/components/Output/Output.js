import React from 'react'
import styles from './Output.module.css'
import OutputCard from './OutputCard'

const Output = (props) => {

    return (
        <section className={styles["output"]}>
            <div className="inner">
                <OutputCard videoId={props.videoId} data={props.data} waiting={props.waiting} />
            </div>
        </section>
    )
}

export default Output
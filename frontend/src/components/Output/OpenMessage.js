import React from "react"
import styles from "./OpenMessage.module.css"
import Loader from "react-loader-spinner"


const OpenMessage = (props) => {
    if (props.isLoading) {
        return (
            <Loader
                type="TailSpin"
                color="#00BFFF"
                height={20}
                width={20}
            />
        )
    }
    return (
        <span><i className={`fas fa-external-link-alt ${styles["external"]}`}></i>Open In New Tab</span>
    )
}

export default OpenMessage
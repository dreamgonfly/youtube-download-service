import React from "react"
import styles from "./OptionsBlock.module.css"

const OptionsBlock = (props) => {
    const options = props.data.Formats.filter((element) => { return element.FormatId !== props.formatId })

    if (options.length === 0) {
        return (
            <div className={styles["output-format-options"]}>
                <div className={styles["option-desc"]}>No other options</div>
            </div>
        )
    }

    const optionClickHandler = (event) => {
        props.onSetFormatId(event.currentTarget.dataset.index)
        props.onSetOptionsOpen(false)
    }

    return (
        <div className={styles["output-format-options"]}>
            {options.map((item, index) => {
                return (
                    <div key={index} data-index={item.FormatId} className={styles["option-desc"]} onClick={optionClickHandler}>
                        {item.Ext} {item.FormatNote} {Math.round(item.Filesize / 1000 / 1000 * 10) / 10} MB
                    </div>
                )
            })}
        </div>
    )
}

export default OptionsBlock
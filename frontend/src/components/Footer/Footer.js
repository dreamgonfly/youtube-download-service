import React from 'react'
import styles from './Footer.module.css'

const Footer = () => {
    return (
        <footer>
            <div className="inner">
                <div className={styles["footer-message"]}>Online Youtube Downloader</div>
                <div className={styles["footer-contact"]}>dreamgonfly@gmail.com</div>
                <div className={styles["footer-copyright"]}>Copyright 2021 © All rights reserved.</div>
            </div>
        </footer>
    )
}

export default Footer
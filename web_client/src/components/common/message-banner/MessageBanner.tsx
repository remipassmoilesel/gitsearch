import * as React from 'react'
import {Component, ReactNode} from 'react'

import './MessageBanner.scss'

export class MessageBanner extends Component<{}, {}> {
    public render(): ReactNode {
        return (
            <div className="gs-message-banner" id="gs-message-banner">&nbsp;</div>
        );
    }
}

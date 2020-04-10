import * as React from 'react'
import {Component, ReactNode} from 'react'

import './Banner.scss'

export class Banner extends Component<{}, {}> {
    public render(): ReactNode {
        return (
            <div className="gs-banner">user@hostname:/path/to/repository</div>
        );
    }
}

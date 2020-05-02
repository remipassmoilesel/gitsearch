import * as React from 'react'
import {Component, ReactNode} from 'react'
import './ControlBar.scss'

interface State {
    processing: boolean
}

export class ControlBar extends Component<{}, State> {

    constructor(props: {}) {
        super(props);
        this.state = {
            processing: false
        }
    }

    public render(): ReactNode {
        const isProcessing = this.state.processing;
        return (
            <div className="gs-control-bar">
                {isProcessing && "Processing ..."}
            </div>
        );
    }

}

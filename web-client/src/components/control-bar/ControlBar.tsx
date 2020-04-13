import * as React from 'react'
import {Component, ReactNode} from 'react'
import {servicesMap} from '../../services/service-factory';
import {faCog, faTrash} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {MessageBanner} from "./message-banner/MessageBanner";
import './ControlBar.scss'

interface State {
    processing: boolean
}

export class ControlBar extends Component<{}, State> {

    private services = servicesMap();

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
                <button onClick={this.buildIndex} disabled={isProcessing} className="btn btn-primary">
                    <FontAwesomeIcon icon={faCog} className="icon"/>Build index
                </button>
                <button onClick={this.cleanIndex} disabled={isProcessing} className="btn btn-secondary">
                    <FontAwesomeIcon icon={faTrash} className="icon"/>Clean index
                </button>
                <MessageBanner/>
            </div>
        );
    }

    private buildIndex = () => {
        this.setState({processing: true});
        this.services.index.buildIndex()
            .then(res => this.services.uiMessages.info(`Build took ${res.TookSeconds} seconds, ${res.Files} files indexed`))
            .finally(() => this.setState({processing: false}));
    };

    private cleanIndex = () => {
        this.setState({processing: true});
        this.services.index.cleanIndex()
            .then(res => this.services.uiMessages.info(`Clean took ${res.TookMillis} ms`))
            .finally(() => this.setState({processing: false}));
    };
}

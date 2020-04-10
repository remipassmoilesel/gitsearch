import * as React from 'react'
import {ChangeEvent, Component, ReactNode} from 'react'
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faSearch} from '@fortawesome/free-solid-svg-icons'

import './SearchBar.scss'

export interface SearchEvent {
    query: string
}

interface Props {
    onSearch: (ev: SearchEvent) => void
}

interface State {
    query: string
}

export class SearchBar extends Component<Props, State> {

    constructor(props: Props) {
        super(props);
        this.state = {
            query: ""
        }
    }

    public render(): ReactNode {
        return (
            <div className="input-group gs-search-bar">
                <input type="text"
                       value={this.state.query}
                       onChange={this.onUserInput}
                       className="form-control"
                       placeholder="Search ..."/>
                <div className="input-group-append">
                    <button className="btn btn-primary" onClick={this.onSearch}>
                        <FontAwesomeIcon icon={faSearch}/>
                    </button>
                </div>
            </div>
        );
    }

    private onUserInput = (ev: ChangeEvent<HTMLInputElement>) => {
        this.setState({query: ev.target.value})
    };

    private onSearch = () => {
        this.props.onSearch({query: ""})
    };
}

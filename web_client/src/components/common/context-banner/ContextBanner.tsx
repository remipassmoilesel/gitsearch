import * as React from 'react'
import {Component, ReactNode} from 'react'
import {RepositoryContext} from "../../../core/repo/RepositoryContext";
import {servicesMap} from "../../../core/service-factory";
import './ContextBanner.scss'

interface State {
    context: RepositoryContext | undefined
}

export class ContextBanner extends Component<{}, State> {

    private services = servicesMap();

    constructor(props: {}) {
        super(props);
        this.state = {
            context: undefined
        }
    }

    public componentDidMount(): void {
        this.services.repository.getContext()
            .then(res => this.setState({context: res}))
    }

    public render(): ReactNode {
        const bannerText = this.getBannerText();
        return (
            <div className="gs-context-banner">GitSearch - {bannerText}</div>
        );
    }

    private getBannerText(): string {
        if (this.state.context) {
            const ctx = this.state.context;
            return `${ctx.Username}@${ctx.Hostname}:${ctx.Path}`
        }
        return "..."
    }
}

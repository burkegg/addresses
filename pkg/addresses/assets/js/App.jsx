class App extends React.Component {
  // a list
  constructor(props) {
    super(props)
    this.state = {
      houses: [],
      version: '',
      searchTerm: '',
    }
  }
  componentDidMount() {

  }

  handleSearchChange = async evt => {
    await this.setState({ searchTerm: evt.target.value })
    if (this.state.searchTerm !== '') {
      await this.handleFetches()
    }
  }

  handleFetches = async () => {
    let search = { Term: this.state.searchTerm}
    let initReq = {
      method: "POST",
      body: JSON.stringify(search)
    }
    const fetchEndpoint = "api/addresses"
    let resp = await fetch(fetchEndpoint, initReq)
    let searchResults = await resp.json()
    await this.setState({ houses: searchResults })
  }

  render() {
    return (
      <React.Fragment>
        <form>
          <label>
            Search for properties:
            <input type={'text'} id={'searchInput'} name={'search'} value={this.state.searchTerm} onChange={this.handleSearchChange}/>
          </label>
        </form>
        <div>
          {
            this.state.houses.map(property => {
              return (
                <div key={property.ID}>{property.Address}</div>
              )
            })
          }
        </div>
      </React.Fragment>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));

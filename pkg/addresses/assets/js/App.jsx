const { Component } = React

const DisplayBox = props => {
  return (
    // whole box
    <a href={props.URL} target={'_blank'}>
      <div style={{display: 'flex', flexDirection: 'row', height: 350, backgroundColor: 'lightgrey', borderStyle: 'solid', borderWidth: '2px', borderColor: 'black', margin: 'auto', marginBottom: '20px'}}>
          {/* Left hand column - house and price */}
          <div className={'image-price'} style={{display: 'flex', flexDirection: 'column'}}>
            <img src={'./house.svg'} alt={"house svg"} style={{height: 250, borderStyle: 'solid', borderColor: 'black', margin: 'auto', marginTop: 30}} zIndex={2}/>
            <p style={{fontSize: '3rem', fontStyle: 'bold', margin: 'auto', justifyText: 'center'}}>
              ${props.price}
            </p>
          </div>
          {/* Right hand column - some info about house */}
          <div style={{display: 'flex', flexDirection: 'column', margin: 'auto'}}>
            <p className={'address-city'} style={{fontSize:'3rem', marginBottom: -10}}>
              {`${props.address},`}
            </p>
            <p className={'address-city'} style={{fontSize:'3rem'}}>
              {props.city}
            </p>
            <p className={'property-type'} style={{fontSize:'2rem'}}>
              {`Property Type: ${props.propertyType}`}
            </p>
          </div>
      </div>
    </a>
  )
}

class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      houses: [],
      version: '',
      searchTerm: '',
    }
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
      <div>
        <form style={{marginTop: 30, marginLeft: 100}}>
          <label>
            Search for properties by street address:
            <input style={{marginTop: 30, marginLeft: 100, marginBottom: 20}} type={'text'} id={'searchInput'} name={'search'} value={this.state.searchTerm} onChange={this.handleSearchChange}/>
          </label>
        </form>
        <div>
          {
            this.state.houses.map(property => {
              return (
                <div style={{width: '60%', marginLeft: 30}} key={property.ID}>
                  <DisplayBox
                    URL={property.URL}
                    address={property.Address}
                    city={property.City}
                    price={property.Price}
                    sqFeet={property.SqFeet}
                    propertyType={property.PropType}
                    />
                </div>
              )
            })
          }
        </div>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));

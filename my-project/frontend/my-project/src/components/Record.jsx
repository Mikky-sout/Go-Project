const Record = ({data}) => {
     const name = "พักโรงแรม"
     const amount = -5000
     return (
       <li className="sub-item">{data} <span>{amount}</span></li>
     )
}

export default Record
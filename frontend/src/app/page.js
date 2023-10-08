"use client"

import Image from 'next/image'
import { useEffect, useState } from 'react'

export default function Home() {
  const BASE_URL = process.env.NEXT_PUBLIC_IS_COMPOSE
    ? ""
    : "http://localhost:3001";

  const columns = [
    { key: "serial_number", label: "Serial #", transform: transformUpperCase },
    { key: "state", label: "Allocation state", transform: transformTitle },
    { key: "software_version", label: "Software version" },
    { key: "product_code", label: "Product code", transform: transformUpperCase }
  ]

  const [komps, setKomps] = useState([])
  const [selectedKomp, setSelectedKomp] = useState(null);

  const fetchKomps = async () => {
    const res = await fetch(`${BASE_URL}/api/komps`)
    const data = await res.json()
    setKomps(data)
  }
  useEffect(() => { fetchKomps() }, []);

  const saveSelectedKompComment = async () => {
    let komp = await patchKomp(selectedKomp.serial_number, { comment: selectedKomp.comment })
    setSelectedKomp(komp);
  }

  const allocateSelectedKomp = async () => {
    let komp = await patchKomp(selectedKomp.serial_number, { state: "allocated" })
    setSelectedKomp(komp);
  }

  const patchKomp = async (serial_number, patch) => {
    const res = await fetch(`${BASE_URL}/api/komps/${serial_number}`, {
      method: 'PATCH',
      body: JSON.stringify(patch),
      headers: { 'Content-type': 'application/json' },
    });

    const data = await res.json();
    await fetchKomps()
    return data;
  }

  const resetSelectedKomp = async () => {
    let komp = await patchKomp(selectedKomp.serial_number, { state: "available" })
    setSelectedKomp(komp)
  }

  return (
    <main className="flex min-h-screen flex-col justify-between p-12">
      <div className="z-10 w-full font-mono lg:flex flex-col">
        <Image src="/komp-logo.svg" alt="Komp logo" className="pb-10" width={100} height={24} priority />
        <div className="flex w-full text-4xl font-semibold pb-12">
          No Isolation - BoVel Komp registry
        </div>

        <div className="w-full flex">
          <div className="w-2/3 pr-20">
            <table className="items-center bg-transparent w-full border-collapse">

              <thead>
                <KompTableHeader columns={columns} />
              </thead>

              <tbody>
                {
                  komps.map(komp =>
                    <KompTableRow
                      key={komp.serial_number}
                      columns={columns}
                      komp={komp}
                      selected={komp.serial_number == selectedKomp?.serial_number}
                      selectKomp={setSelectedKomp}
                    />
                  )
                }
              </tbody>

            </table>
          </div>
          <div className="w-1/3 right-12 border border-solid border-t-0 border-b-0 border-r-0">

            {
              !selectedKomp &&
              <div className="w-full h-full flex items-center break-normal justify-center">Select a Komp</div>
            }

            {
              selectedKomp &&
              <div className="w-full flex flex-col pl-12">
                <div className="text-2xl">Komp details</div>
                <div className="pl-2 text-lg w-full space-y-4 pt-5">
                  <div>Serial number: {transformUpperCase(selectedKomp.serial_number)}</div>
                  <div>Allocation status: {transformTitle(selectedKomp.state)}</div>
                  <div>Software version: {selectedKomp.software_version}</div>
                  <div>Product code: {transformUpperCase(selectedKomp.product_code)}</div>
                  <div>MAC address: {transformUpperCase(selectedKomp.mac_address)}</div>

                  {selectedKomp.attributes.length ?
                    <div>
                      <div>Attributes:</div>
                      <div className="pl-2">
                        {
                          (selectedKomp.attributes || [])
                            .map(attribute => (<div key={attribute.name}>{attribute.name} = {attribute.value}</div>))
                        }
                      </div>
                    </div> : null
                  }

                  <div className="">
                    <div>Comments:</div>
                    <form onSubmit={(e) => { saveSelectedKompComment(); e.preventDefault(); }}>
                      <div>
                        <textarea
                          className="bg-stone-100 p-2 w-full"
                          rows={7}
                          value={selectedKomp.comment || ""}
                          onChange={(e) => { setSelectedKomp({ ...selectedKomp, comment: e.target.value }); }}
                        />
                      </div>

                      <div className="flex align-right pt-2">
                        <button className="hover:bg-stone-200 bg-stone-100 border border-solid p-2" type="submit">Save comment</button>
                      </div>
                    </form>
                  </div>

                  {
                    selectedKomp.state == "allocated" ?
                      <form onSubmit={(e) => { resetSelectedKomp(); e.preventDefault(); }}>
                        <button className="hover:bg-red-700 bg-red-800 text-stone-100 p-2" type="submit">Reset</button>
                      </form> :
                      <form onSubmit={(e) => { allocateSelectedKomp(); e.preventDefault(); }}>
                        <button className="hover:bg-stone-200 bg-stone-100 border border-solid p-2" type="submit">Allocate</button>
                      </form>
                  }
                </div>
              </div>
            }
          </div>
        </div>
      </div>
    </main >
  )
}

function KompTableHeader(props) {
  const { columns } = props;
  return (
    <tr>
      {
        columns.map(column => {
          return (
            <th
              key={column.label}
              className="px-6 bg-stone-100 align-middle border border-solid py-3 border-l-0 border-r-0 whitespace-nowrap font-semibold text-lg text-left">
              {column.label}
            </th>
          )
        })
      }
    </tr>
  )
}


function KompTableRow(props) {
  const { columns, komp, selectKomp, selected } = props;
  return (
    <tr
      className={`cursor-pointer ${selected ? 'bg-gray-100' : ''} hover:bg-gray-100`}
      onClick={() => selectKomp(komp)}>
      {
        columns.map(column => {
          const value = column.transform ? column.transform(komp[column.key]) : komp[column.key];
          return (
            <td
              key={`${komp.serial_number}[${column.key}]`}
              className="border-t-0 px-6 align-middle border-l-0 border-r-0 whitespace-nowrap p-4 text-left text-base font-extralight">
              {value}
            </td>
          )
        })}
    </tr>
  )
}

function transformUpperCase(value) {
  if (!value || typeof (value) != "string") {
    return value
  }
  return value.toUpperCase()
}

function transformTitle(value) {
  if (!value || typeof (value) != "string") {
    return value
  }
  return value[0].toUpperCase() + value.slice(1)
}
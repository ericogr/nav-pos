var vectorTileOptions = {
    vectorTileLayerStyles: {
        water: {
            fill: true,
            weight: 1,
            fillColor: '#06cccc',
            color: '#06cccc',
            fillOpacity: 0.2,
            opacity: 0.4,
        },
        admin: {
            weight: 1,
            fillColor: 'pink',
            color: 'pink',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        waterway: {
            weight: 1,
            fillColor: '#2375e0',
            color: '#2375e0',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        landcover: {
            fill: true,
            weight: 1,
            fillColor: '#53e033',
            color: '#53e033',
            fillOpacity: 0.2,
            opacity: 0.4,
        },
        landuse: {
            fill: true,
            weight: 1,
            fillColor: '#e5b404',
            color: '#e5b404',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        mountain_peak: {
            fill: false,
            stroke: false,
            fillColor: "#ff6347",
            color: "#b23c2f",
            weight: 2
        },
        park: {
            fill: true,
            weight: 1,
            fillColor: '#84ea5b',
            color: '#84ea5b',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        boundary: {
            stroke: false,
            fill: false,
            weight: 1,
            fillColor: '#c545d3',
            color: '#c545d3',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        aeroway: {
            weight: 1,
            fillColor: '#51aeb5',
            color: '#51aeb5',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        transportation: (properties, zoom) => {
            const weight = zoom <= 10 ? 0.6 : zoom <= 14 ? 1 : 1.4;
            return {
                stroke: true,
                color: "#444",
                weight: weight,
                opacity: 0.9
            };
        },
        building: {
            fill: true,
            weight: 1,
            fillColor: '#2b2b2b',
            color: '#2b2b2b',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        water_name: {
            stroke: false,
            fill: false,
            weight: 1,
            fillColor: '#022c5b',
            color: '#022c5b',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        transportation_name: {
            weight: 1,
            fillColor: '#bc6b38',
            color: '#bc6b38',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        place: {
            stroke: false,
            fill: false,
            weight: 1,
            fillColor: '#f20e93',
            color: '#f20e93',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        housenumber: {
            stroke: false,
            fill: false,
            weight: 1,
            fillColor: '#ef4c8b',
            color: '#ef4c8b',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        poi: {
            stroke: false,
            fill: false,
            weight: 1,
            fillColor: '#3bb50a',
            color: '#3bb50a',
            fillOpacity: 0.2,
            opacity: 0.4
        },
        aerodrome_label: {
            stroke: false,
            fill: false,
            textSize: "14px",
            color: "#ff6347",
            fontWeight: "bold"
        }
    },
    maxZoom: 14
};

const initialCoordinates = {
    center: [(-23.607475903452663 + -24.270974395948752) / 2,
    (-46.48307030043036 + -46.97940265774286) / 2],
    bounds: [
        [-24.270974395948752, -46.97940265774286],
        [-23.607475903452663, -46.48307030043036]
    ]
};
const aircraftMarkers = {};
const HEADER_SESSION_ID = 'X-Session-ID'
const globalMap = L.map('map').fitBounds(initialCoordinates.bounds);

let telemetryMarker = null;
let debounceTimer;
let centerOnTelemetry = true;
let lastTelemetryData = null;
let sessionInfo = {};

setTimeout(() => {
    updateTelemetry();
    setInterval(updateTelemetry, 1000);
}, 1000)
setTimeout(() => {
    updateAircrafts();
    setInterval(updateAircrafts, 30000);
}, 5000)
updateSessionId();

globalMap.on('moveend', debouncedUpdateAircraftData);
globalMap.on('zoomend', debouncedUpdateAircraftData);

document.getElementById('centerButton').addEventListener('click', function () {
    if (lastTelemetryData && lastTelemetryData.lat && lastTelemetryData.lon) {
        globalMap.setView([lastTelemetryData.lat, lastTelemetryData.lon], globalMap.getZoom());
    }
});

function applyTileMap(map, name) {
    if (name === 'mbtiles') {
        L.vectorGrid.protobuf('/api/v1/tile/{z}/{x}/{y}.img', vectorTileOptions).addTo(map);
        console.info('Using mbtiles');
    }
    else {
        L.tileLayer('/api/v1/tile/{z}/{x}/{y}.img').addTo(map);
    }
}

function createAircraftIcon(heading, color) {
    return L.divIcon({
        className: 'aircraft-marker',
        html: `
        <svg width="32" height="32" viewBox="0 0 1280 1280"
                xmlns="http://www.w3.org/2000/svg"
                style="transform: rotate(${heading}deg);">
            <g transform="translate(0,1280) scale(0.1,-0.1)">
                <path d="M6245 12784 c-340 -71 -562 -380 -667 -924 -22 -115 -22 -120 -25
                -1507 l-4 -1391 -27 -10 c-317 -114 -2812 -1003 -4477 -1597 -473 -168 -900
                -321 -950 -339 l-90 -33 -3 -716 c-1 -474 1 -717 8 -717 5 0 1253 191 2772
                425 1519 234 2762 424 2764 422 1 -2 263 -3932 267 -4021 l2 -50 -1200 -736
                -1200 -736 -3 -428 -2 -428 22 5 c13 3 84 17 158 32 217 42 1219 237 1670 325
                228 45 505 99 615 120 110 21 273 53 362 71 l163 31 162 -31 c90 -18 253 -50
                363 -71 110 -21 387 -75 615 -120 451 -88 1453 -283 1670 -325 74 -15 145 -29
                158 -32 l22 -5 -2 428 -3 428 -1200 736 -1200 736 2 50 c4 89 266 4019 267
                4021 2 2 1245 -188 2764 -422 1519 -234 2767 -425 2772 -425 7 0 9 243 8 717
                l-3 717 -2770 1066 -2770 1065 -6 1315 c-5 1302 -5 1316 -27 1430 -89 464
                -258 748 -517 871 -139 66 -306 85 -460 53z"
                fill="${color}" stroke="#000000" stroke-width="200"/>
            </g>
        </svg>`,
        iconSize: [32, 32],
        iconAnchor: [16, 16]
    });
}

function createBlueAircraftIcon(heading) {
    return createAircraftIcon(heading, '#3388ff');
}

function createRedAircraftMeIcon(heading) {
    return createAircraftIcon(heading, '#ff3333');
}

function updateSessionId() {
    fetch('/api/v1/session')
        .then(response => {
            if (!response.ok) {
                throw new Error(`Session error: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            sessionInfo = data;
            console.info('Session ID:', sessionInfo.sessionId);

            applyTileMap(globalMap, sessionInfo.tileMapServiceName);
        })
        .catch(error => {
            console.error('Error fetching session ID:', error);
        });
}

function updateAircrafts() {
    const bounds = globalMap.getBounds();
    const bbox = {
        min_latitude: bounds.getSouth(),
        min_longitude: bounds.getWest(),
        max_latitude: bounds.getNorth(),
        max_longitude: bounds.getEast()
    };

    document.querySelector('.area-info').textContent =
        `Area: ${bbox.min_latitude.toFixed(4)},${bbox.min_longitude.toFixed(4)} to ${bbox.max_latitude.toFixed(4)},${bbox.max_longitude.toFixed(4)}`;

    fetch('/api/v1/radar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            [HEADER_SESSION_ID]: sessionInfo.sessionId
        },
        body: JSON.stringify(bbox)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Aircrafts error: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            const now = new Date();
            document.querySelector('.last-update').textContent = `Last update: ${now.toLocaleTimeString()}`;

            if (data == null || data.length == 0) {
                document.querySelector('.aircraft-count').textContent = 'No visible aircrafts';
                return;
            }

            document.querySelector('.aircraft-count').textContent = `Visible aircrafts: ${data.length}`;

            const updatedAircrafts = new Set();

            data.forEach(aircraft => {
                if (!aircraft.latitude || !aircraft.longitude) return;

                updatedAircrafts.add(aircraft.icao24);

                if (aircraftMarkers[aircraft.icao24]) {
                    aircraftMarkers[aircraft.icao24].setLatLng([aircraft.latitude, aircraft.longitude]);

                    if (aircraft.true_track !== null) {
                        aircraftMarkers[aircraft.icao24].setIcon(createBlueAircraftIcon(aircraft.true_track));
                    }

                    aircraftMarkers[aircraft.icao24].getPopup().setContent(createPopupContent(aircraft));
                } else {
                    const marker = L.marker([aircraft.latitude, aircraft.longitude], {
                        icon: createBlueAircraftIcon(aircraft.true_track || 0)
                    }).addTo(globalMap);

                    marker.bindPopup(createPopupContent(aircraft));

                    aircraftMarkers[aircraft.icao24] = marker;
                }
            });

            Object.keys(aircraftMarkers).forEach(icao24 => {
                if (!updatedAircrafts.has(icao24)) {
                    globalMap.removeLayer(aircraftMarkers[icao24]);
                    delete aircraftMarkers[icao24];
                }
            });
        })
        .catch(error => {
            document.querySelector('.aircraft-count').textContent = error.message;
        });
}

function updateTelemetry() {
    fetch('/api/v1/telemetry', {
        headers: {
            [HEADER_SESSION_ID]: sessionInfo.sessionId
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Telemetry error: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (!data || !data.lat || !data.lon) {
                telemetryMessage("no data");
                return;
            }

            lastTelemetryData = data;

            document.querySelector('.telemetry-info').style.display = 'block';
            document.getElementById('centerButton').style.display = 'block';

            document.querySelector('.telemetry-title').textContent = `MAVLink Telemetry Data`;
            document.querySelector('.telemetry-position').textContent = `Position: ${data.lat.toFixed(6)}, ${data.lon.toFixed(6)}`;
            document.querySelector('.telemetry-altitude').textContent = `Altitude: ${data.alt.toFixed(1)} m`;
            document.querySelector('.telemetry-heading').textContent = `Heading: ${data.hdg.toFixed(1)}°`;
            document.querySelector('.telemetry-speed').textContent = `Speed: ${data.spd.toFixed(1)} km/h`;
            document.querySelector('.telemetry-satellites').textContent = `Satellites: ${data.sat}`;
            document.querySelector('.telemetry-battery').textContent = `Battery: ${data.bat.toFixed(1)}V / ${data.rem}%`;
            telemetryMessage("valid");

            if (telemetryMarker) {
                telemetryMarker.setLatLng([data.lat, data.lon]);
                telemetryMarker.setIcon(createRedAircraftMeIcon(data.hdg));
            } else {
                telemetryMarker = L.marker([data.lat, data.lon], { icon: createRedAircraftMeIcon(data.hdg) }).addTo(globalMap);
                telemetryMarker.bindPopup(createTelemetryPopupContent(data));
            }

            if (centerOnTelemetry && data.valid) {
                globalMap.setView([data.lat, data.lon], globalMap.getZoom());
                centerOnTelemetry = false;
            }
        })
        .catch(error => {
            document.querySelector('.telemetry-title').textContent = error.message;
            telemetryMessage("error fetching telemetry data");
        });
}

function createPopupContent(aircraft) {
    return `
                <div class="aircraft-info">
                    <strong>Callsign:</strong> ${aircraft.callsign || 'N/A'}<br>
                    <strong>ICAO24:</strong> ${aircraft.icao24}<br>
                    <strong>Country Origin:</strong> ${aircraft.origin_country}<br>
                    <strong>Altitude:</strong> ${aircraft.baro_altitude ? Math.round(aircraft.baro_altitude) + ' m' : 'N/A'}<br>
                    <strong>Speed:</strong> ${aircraft.velocity ? Math.round(aircraft.velocity * 3.6) + ' km/h' : 'N/A'}<br>
                    <strong>Heading:</strong> ${aircraft.true_track ? Math.round(aircraft.true_track) + '°' : 'N/A'}<br>
                    <strong>On Ground:</strong> ${aircraft.on_ground ? 'Yes' : 'No'}
                </div>
            `;
}

function createTelemetryPopupContent(telemetry) {
    return `
                <div class="aircraft-info">
                    <strong style="color: #d32f2f;">MAVLinkTelemetry</strong><br>
                    <strong>Latitude:</strong> ${telemetry.lat.toFixed(6)}<br>
                    <strong>Longitude:</strong> ${telemetry.lon.toFixed(6)}<br>
                    <strong>Altitude:</strong> ${telemetry.alt.toFixed(1)} m<br>
                    <strong>Speed:</strong> ${telemetry.spd.toFixed(1)} km/h<br>
                    <strong>Heading:</strong> ${telemetry.hdg.toFixed(1)}°<br>
                    <strong>Satelites:</strong> ${telemetry.sat}<br>
                    <strong>Battery:</strong> ${telemetry.bat.toFixed(1)}V / ${telemetry.rem}%<br>
                    <strong>Current:</strong> ${telemetry.cur.toFixed(2)}A<br>
                    <strong>Consumption:</strong> ${telemetry.mah}mAh
                </div>
            `;
}

function telemetryMessage(message) {
    document.querySelector('.telemetry-valid').textContent = `Message: ${message}`;
}

function debouncedUpdateAircraftData() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(updateAircrafts, 500);
}
import http from "k6/http";

export default function() {
    let response = http.post(__ENV.API_URL + "/transfers",  JSON.stringify({
        amount: 1,
        debtor_id : 'f2f0e0d1-e37e-45c3-ad06-e6c2a66544fc',
        beneficiary_id : '089557bc-ddf2-4ec5-8077-d8bf09fe3ddc' 
    }));

    http.setResponseCallback(http.expectedStatuses(201));
};
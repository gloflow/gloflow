/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

// ///<reference path="../../d/jquery.d.ts" />

//-------------------------------------------------
export function add_query_param(p_url_str :string, key :string, value :string) {
    let url = new URL(p_url_str);
    url.searchParams.set(key, value); // Add or update parameter
    return url.toString();
}

//-------------------------------------------------
export function get_current_host() :string {
	const domain_str   = window.location.hostname;
	const protocol_str = window.location.protocol;
	const host_str :string = `${protocol_str}//${domain_str}`;

	console.log("gf_host", host_str);
	return host_str;
}

//-------------------------------------------------
export function click_outside(p_element :any, p_on_click_fun :any) {
    $(document).on("click", (p_e :any)=>{
            
        // check if a click is not on the element or a child of the element
        if (!p_element.is(p_e.target) && $(p_element).has(p_e.target).length === 0) {
            p_on_click_fun()
        }
    });
}

//--------------------------------------------------------
// MEASURE_DIV_DIMENSIONS

export function measureDivDimensions(pTargetElement :HTMLElement,
    pMeasurementParentElement :HTMLElement) {

    const isAlreadyInDOM = document.body.contains(pTargetElement);

    // IN_DOM
    if (isAlreadyInDOM) {
        return {
            width:  pTargetElement.offsetWidth,
            height: pTargetElement.offsetHeight
        };
    }

    // NOT_IN_DOM
    else {
        
        /*
        // container for off-screen measurement
        const offScreenContainer = document.createElement('div');
        offScreenContainer.style.position = 'absolute';
        offScreenContainer.style.left = '-9999px';
        offScreenContainer.style.top = '-9999px';
        document.body.appendChild(offScreenContainer);
        */

        // add the div to the off-screen container
        pMeasurementParentElement.appendChild(pTargetElement);

        const originalVisibility = pTargetElement.style.visibility;
        pTargetElement.style.visibility = 'hidden';

        // measure the div
        const dimensionsMap = {
            width:  pTargetElement.offsetWidth,
            height: pTargetElement.offsetHeight
        };

        // remove the off-screen container
        // document.body.removeChild(offScreenContainer);

        pTargetElement.style.visibility = originalVisibility;
        pMeasurementParentElement.removeChild(pTargetElement);

        return dimensionsMap;
    }
}
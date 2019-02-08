/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

///<reference path="../../../d/jquery.d.ts" />

//--------------------------------------------------------
export function init(p_register_user_email_fun, p_log_fun) {
	p_log_fun('FUN_ENTER','gf_email_registration.init()');

	const register_email_form = $(`
		<div id='register_email_form'>
			<form>
				<!-- user email -->
				<div>
					<input type       ="text"
						   size       ="60" 
						   id         ="user_email_input" 
						   placeholder="your email..."></input>
				</div>
				<div id='submit_register_email_form_button' class='button'>
					<div class='button_title'>send</div>
				</div>
			</form>
		</div>`);
	$('#register').append(register_email_form);

	const submit_register_email_form_button = $(register_email_form).find('#submit_register_email_form_button');
	
	//this button reveals the email registration form
	var email_form_visible_bool :boolean = false;
	$('#register #register_email_button').click((p_event)=>{
		
		if (email_form_visible_bool == false) {
			$(register_email_form).css('opacity','1.0');
			email_form_visible_bool = true;

			//layout_email_form(p_log_fun);
		}
		else {
			$(register_email_form).css('opacity','0.0');
			email_form_visible_bool = false;
		}
	});

	$(submit_register_email_form_button).click((p_event)=>{
			const inputed_email_str :string = $(register_email_form).find('#user_email_input').val();
			register_user_email(inputed_email_str);
		});
	
	//--------------------------------------------------------
	function register_user_email(p_inputed_email_str :string) {
		p_log_fun('FUN_ENTER', 'gf_email_registration.init().register_user_email()');

		p_register_user_email_fun(p_inputed_email_str,

			//p_onComplete_fun
			(p_status_str :string,
			p_msg_str     :string)=>{
				console.assert(p_status_str == 'success' || p_status_str == 'error');

				p_log_fun('INFO', 'email registration DONE');
				p_log_fun('INFO', 'p_status_str:$p_status_str');

				switch(p_status_str) {
					case 'success':
						$(submit_register_email_form_button).find('.button_title').text('success');
						$(submit_register_email_form_button).css('background-color', 'rgb(80, 173, 36)');
						break;
					case 'error':
						$(submit_register_email_form_button).find('.button_title').text('error');
						$(submit_register_email_form_button).css('background-color', 'rgb(255, 10, 0)');

						const error_msg = $(`
							<div class="button_title">$p_msg_str</div>
						`);
						$(submit_register_email_form_button).append(error_msg);
						break;
				}
			},
			p_log_fun);
	}
	//--------------------------------------------------------
}
/*//--------------------------------------------------------
function layout_email_button(p_log_fun) {
	p_log_fun('FUN_ENTER','gf_email_registration.layout_email_button()');

	final DivElement btn_element = query('#register_email_button');

	if (window.innerWidth < 770) {

		if (btn_element.dataset.containsKey('minimized_bool') &&
			btn_element.dataset['minimized_bool'] == 'false') {
			final int original_width_int = int.parse(btn_element.getComputedStyle().width.replaceAll('px',''));
			btn_element.dataset['original_width_int'] = original_width_int.toString();
			btn_element.dataset['minimized_bool']     = 'true';
			btn_element.style.width                   = '71px';
			btn_element.query('.button_title').text   = 'i';
		}
	}
	else {
		btn_element.dataset['minimized_bool']   = 'false';
		btn_element.style.width                 = '''${btn_element.dataset['original_width_int']}px''';
		btn_element.query('.button_title').text = 'Get an Invite';
	}
}
//--------------------------------------------------------
function layout_email_form(p_log_fun) {
	p_log_fun('FUN_ENTER','gf_email_registration.layout_email_form()');

	final DivElement register_email_form = $('#register_email_form');

	const button_x_int     :number = int.parse($('#register_email_button').getComputedStyle().left.replaceAll('px',''));
	const button_width_int :number = int.parse($('#register_email_button').getComputedStyle().width.replaceAll('px',''));
	const form_width_int   :number = int.parse(register_email_form.getComputedStyle().width.replaceAll('px',''));
	const form_x_int       :number;

	//check that the form fits into the window (that a part of it will not be obscured).
	//if it is then right align it so that the user can view it in full
	if ((button_x_int + form_width_int) > window.innerWidth) {
		form_x_int = (button_x_int + button_width_int) - form_width_int;
	}
	else {
		form_x_int = button_x_int;
	}
	register_email_form.style.left = '${form_x_int}px';
}*/
/*
 * DAInstruction.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DAInstruction.h"

namespace dra {

	DAInstruction::DAInstruction() {
		state = Kind::other;

		parent = nullptr;

	}

	DAInstruction::~DAInstruction() = default;

} /* namespace dra */

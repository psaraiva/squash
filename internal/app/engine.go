package app

func (p *Squash) Update() {
	if p.State != StatePlaying {
		return
	}

	p.calcNextLevel()
	p.calcMoveBall()
	p.calcCollisionWithWalls()
	p.calcCollisionPaddle()
	p.calcLostLive()
}

func (p *Squash) CalcMovePaddle(axisY float64) {
	if p.State != StatePlaying {
		return
	}

	newY := axisY - (p.PaddleH / 2)
	if newY < 0 {
		p.PaddleY = 0
	} else if newY > p.Height-p.PaddleH {
		p.PaddleY = p.Height - p.PaddleH
	} else {
		p.PaddleY = newY
	}
}

func (p *Squash) SetPaddlePosition(y float64) {
	if y < 0 {
		p.PaddleY = 0
	} else if y > p.Height-p.PaddleH {
		p.PaddleY = p.Height - p.PaddleH
	} else {
		p.PaddleY = y
	}
}

func (p *Squash) calcCollisionWithWalls() {
	// Top and Bottom walls
	if p.BallY <= 0 || p.BallY >= p.Height-p.BallSize {
		p.BallDY = -p.BallDY
	}

	// Right wall
	if p.BallX >= p.Width-p.BallSize {
		p.BallDX = -p.BallDX
		p.BallX = p.Width - p.BallSize - 1
	}
}

func (p *Squash) calcCollisionPaddle() {
	impactZone := p.PaddleX + p.PaddleW
	if p.BallX+p.BallSize >= p.PaddleX && p.BallX <= impactZone {
		if p.BallY+p.BallSize >= p.PaddleY && p.BallY <= p.PaddleY+p.PaddleH {
			p.BallDX = -p.BallDX
			p.BallX = impactZone + 1
			p.Score += PointsPerCollision
		}
	}
}

func (p *Squash) calcLostLive() {
	if p.BallX+p.BallSize <= 0 {
		p.Lives--
		if p.Lives <= 0 {
			p.State = StateGameOver
		} else {
			p.respawnBall()
		}
	}
}

func (p *Squash) calcNextLevel() {
	currentLevel := p.Score / PointsPerLevel
	if currentLevel > p.LastLevel {

		p.LastLevel = currentLevel
		increment := BaseSpeedBall * p.SpeedIncrement

		if p.BallDX > 0 {
			p.BallDX += increment
		} else {
			p.BallDX -= increment
		}

		if p.BallDY > 0 {
			p.BallDY += increment
		} else {
			p.BallDY -= increment
		}
	}
}

func (p *Squash) calcMoveBall() {
	p.BallX += p.BallDX * p.DeltaTime
	p.BallY += p.BallDY * p.DeltaTime
}

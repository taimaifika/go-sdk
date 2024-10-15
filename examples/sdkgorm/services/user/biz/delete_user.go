package biz

import "context"

func (b *biz) DeleteUserById(ctx context.Context, id int) error {
	if err := b.userRepo.DeleteUserById(ctx, id); err != nil {
		return err
	}

	return nil
}

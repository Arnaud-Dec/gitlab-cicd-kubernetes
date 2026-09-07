package e2e_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"e2e/apiclient"
)

var _ = Describe("Todos API", func() {
	Describe("creating a todo", func() {
		It("returns 201 and the created todo", func(ctx context.Context) {
			resp, err := api.CreateTodo(ctx, "Buy milk", "2% please")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			todo, err := resp.DecodeTodo()
			Expect(err).NotTo(HaveOccurred())
			Expect(todo.ID).To(BeNumerically(">", 0))
			Expect(todo.Title).To(Equal("Buy milk"))
			Expect(todo.Description).To(Equal("2% please"))
			Expect(todo.Completed).To(BeFalse())

			DeferCleanup(func(ctx context.Context) {
				_, _ = api.DeleteTodo(ctx, todo.ID)
			})
		})

		It("rejects a missing title with 400", func(ctx context.Context) {
			resp, err := api.CreateTodo(ctx, "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			body, err := resp.DecodeError()
			Expect(err).NotTo(HaveOccurred())
			Expect(body.Error).NotTo(BeEmpty())
		})
	})

	Describe("the lifecycle of a single todo", func() {
		var todo apiclient.Todo

		BeforeEach(func(ctx context.Context) {
			resp, err := api.CreateTodo(ctx, "Lifecycle todo", "created by the e2e suite")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			todo, err = resp.DecodeTodo()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func(ctx context.Context) {
			// Best-effort: several specs below delete the todo themselves,
			// so a second delete here is expected to 404 and that's fine.
			_, _ = api.DeleteTodo(ctx, todo.ID)
		})

		It("appears in the full list", func(ctx context.Context) {
			resp, err := api.ListTodos(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			todos, err := resp.DecodeTodos()
			Expect(err).NotTo(HaveOccurred())

			ids := make([]uint, len(todos))
			for i, t := range todos {
				ids[i] = t.ID
			}
			Expect(ids).To(ContainElement(todo.ID))
		})

		It("can be fetched by id", func(ctx context.Context) {
			resp, err := api.GetTodo(ctx, todo.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			fetched, err := resp.DecodeTodo()
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.ID).To(Equal(todo.ID))
			Expect(fetched.Title).To(Equal(todo.Title))
		})

		It("can be updated", func(ctx context.Context) {
			resp, err := api.UpdateTodo(ctx, todo.ID, "Updated title", "updated description", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			updated, err := resp.DecodeTodo()
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Title).To(Equal("Updated title"))
			Expect(updated.Description).To(Equal("updated description"))
			Expect(updated.Completed).To(BeTrue())
		})

		It("can be toggled and toggled back", func(ctx context.Context) {
			resp, err := api.ToggleTodo(ctx, todo.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			toggled, err := resp.DecodeTodo()
			Expect(err).NotTo(HaveOccurred())
			Expect(toggled.Completed).To(BeTrue())

			resp, err = api.ToggleTodo(ctx, todo.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			toggledBack, err := resp.DecodeTodo()
			Expect(err).NotTo(HaveOccurred())
			Expect(toggledBack.Completed).To(BeFalse())
		})

		It("returns 404 for it after being deleted", func(ctx context.Context) {
			resp, err := api.DeleteTodo(ctx, todo.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			getResp, err := api.GetTodo(ctx, todo.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Describe("operating on a todo that doesn't exist", func() {
		const missingID = uint(999999999)

		It("get returns 404", func(ctx context.Context) {
			resp, err := api.GetTodo(ctx, missingID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("update returns 404", func(ctx context.Context) {
			resp, err := api.UpdateTodo(ctx, missingID, "x", "", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("toggle returns 404", func(ctx context.Context) {
			resp, err := api.ToggleTodo(ctx, missingID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("delete returns 404", func(ctx context.Context) {
			resp, err := api.DeleteTodo(ctx, missingID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
